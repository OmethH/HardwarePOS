use std::sync::Mutex;
use tauri::{Manager, RunEvent};
use tauri_plugin_shell::process::{CommandChild, CommandEvent};
use tauri_plugin_shell::ShellExt;

struct BackendProcess(Mutex<Option<CommandChild>>);

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  tauri::Builder::default()
    .plugin(tauri_plugin_shell::init())
    .manage(BackendProcess(Mutex::new(None)))
    .setup(|app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }

      let app_data_dir = app.path().app_data_dir().expect("no app data dir");
      let db_path = app_data_dir.join("data").join("hardwarepos.db");
      let migrations_path = app
        .path()
        .resource_dir()
        .expect("no resource dir")
        .join("migrations");

      let sidecar = app
        .shell()
        .sidecar("hardwarepos-backend")
        .expect("failed to create backend sidecar command")
        .env("DB_PATH", db_path.to_string_lossy().to_string())
        .env("MIGRATIONS_PATH", migrations_path.to_string_lossy().to_string())
        .env("PORT", "8971");

      let (mut rx, child) = sidecar.spawn().expect("failed to spawn backend sidecar");
      app.state::<BackendProcess>().0.lock().unwrap().replace(child);

      tauri::async_runtime::spawn(async move {
        while let Some(event) = rx.recv().await {
          match event {
            CommandEvent::Stdout(line) => log::info!("[backend] {}", String::from_utf8_lossy(&line)),
            CommandEvent::Stderr(line) => log::error!("[backend] {}", String::from_utf8_lossy(&line)),
            _ => {}
          }
        }
      });

      Ok(())
    })
    .build(tauri::generate_context!())
    .expect("error while building tauri application")
    .run(|app_handle, event| {
      // The backend sidecar is a separate OS process — Tauri doesn't kill it
      // for us, so it must be stopped explicitly or it lingers after the
      // window closes.
      if let RunEvent::ExitRequested { .. } = event {
        if let Some(child) = app_handle.state::<BackendProcess>().0.lock().unwrap().take() {
          let _ = child.kill();
        }
      }
    });
}
