echo "File changed. Rebuilding..."
while read event; do
  if [[ $event == *"WRITE"* ]] && [[ $event == *".go"* ]]; then
    echo "File changed. Rebuilding..."
  fi
done
