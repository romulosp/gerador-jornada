"""Verifica primeiro uso e preferência no EXE Windows; configurações são fixtures isoladas em build."""

import json
import os
from pathlib import Path
import subprocess
import time

from test_splash_native import capture_window, startup_records


def main():
    """Exercita primeira abertura, reabertura padrão, opção marcada e desmarcada em processos reais."""
    root = Path(__file__).resolve().parents[1]
    evidence = root / "build" / "036-v11-policy"
    evidence.mkdir(parents=True, exist_ok=True)
    # Cada execução da ferramenta usa configuração inédita, sem apagar preferências do usuário.
    config = evidence / f"config-{time.time_ns()}"
    config.mkdir()
    preference_file = config / "startup-preferences.json"
    environment = os.environ.copy()
    environment["JORNADA_CONFIG_DIR"] = str(config)
    executable = root / "dist" / "windows-amd64" / "JornadaStudio.exe"
    # A interação do Check e suas falhas são verificadas em TestLauncherSplashPreferencePersistence.
    for name, should_show, choice in [("first", True, None), ("default-reopen", False, None), ("enabled", True, True), ("disabled", False, False)]:
        if choice is not None:
            preference = json.loads(preference_file.read_text(encoding="utf-8"))
            preference["show_on_startup"] = choice
            preference_file.write_text(json.dumps(preference), encoding="utf-8")
        began = time.monotonic()
        with (evidence / f"{name}-process.log").open("w", encoding="utf-8") as output:
            process = subprocess.Popen([str(executable), "-smoke-screenshot", str(evidence / f"{name}-launcher.png")], env=environment, stdout=output, stderr=subprocess.STDOUT)
            try:
                if should_show:
                    time.sleep(2)
                    capture_window(process.pid, evidence / f"{name}-video-early.png")
                    time.sleep(6)
                    capture_window(process.pid, evidence / f"{name}-video-late.png")
                process.wait(timeout=25 if should_show else 10)
                assert process.returncode == 0, f"{name}: processo encerrou com falha"
                records = startup_records(process.pid)
                events = [record["event"] for record in records]
                assert events.index("startup_resources_ready") < events.index("startup_ready")
                if should_show:
                    assert next(record for record in records if record["event"] == "splash_open")["fullscreen"] is False
                    assert events.index("splash_video_started") < events.index("splash_video_ended") < events.index("startup_ready")
                    ended = next(record for record in records if record["event"] == "splash_video_ended")
                    assert ended["elapsed_ms"] >= 9500, "Reprodução cortada antes do fim"
                    assert ended["audio_track"], "Faixa de áudio não selecionada"
                else:
                    assert "splash_skipped" in events and "splash_open" not in events and "splash_video_started" not in events
                saved = json.loads(preference_file.read_text(encoding="utf-8"))
                assert saved == {"presented": True, "show_on_startup": choice is True}
                (evidence / f"{name}-events.jsonl").write_text("\n".join(json.dumps(record, ensure_ascii=False) for record in records)+"\n", encoding="utf-8")
                print(f"PASS {name}: splash={should_show}, preferência={saved}, duração={time.monotonic()-began:.2f}s", flush=True)
            finally:
                if process.poll() is None:
                    process.terminate()
                    process.wait(timeout=5)


if __name__ == "__main__":
    main()
