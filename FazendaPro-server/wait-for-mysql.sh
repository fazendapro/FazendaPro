until nc -z mysql 3306; do
  echo "Waiting for MySQL to be ready..."
  sleep 1
done
echo "MySQL is up - starting app"
exec "$@"