Start-Sleep -Seconds 5
Start-Process "Chrome" "http://localhost:8000/goradd/Test.g?all=1 --headless --remote-debugging-port=9222"