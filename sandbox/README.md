# Sandbox Runner (opsional)

Mekanisme menjalankan backend Golang sebagai mini-service di sandbox SPKD:
`mini-services/medcredix-api/start.sh` membangun (bila perlu) lalu menjalankan binary di port 8081.
Frontend Vue 3 + Vite berjalan di port 3000 (`npm run dev` di `frontend/`) dengan proxy `/api` ke backend.
