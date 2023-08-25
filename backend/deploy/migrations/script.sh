for file in $(ls *.up.sql | sort -n); do
    psql -U postgres -h localhost -d qaztrade -f "$file"
done
