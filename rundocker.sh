echo "Running your POSTGRES with docker..."
docker-compose up -d
echo "Everything is up!"
echo "Processes:"
ps aux | grep docker-proxy
grep -w 'postgres' /etc/services
