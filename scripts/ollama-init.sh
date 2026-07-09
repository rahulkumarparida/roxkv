#!/bin/sh
set -eu

MODEL="${OLLAMA_MODEL:-llama3.1:latest}"

ollama serve &
OLLAMA_PID=$!

cleanup() {
  kill "${OLLAMA_PID}" 2>/dev/null || true
  wait "${OLLAMA_PID}" 2>/dev/null || true
}

trap cleanup INT TERM

until ollama list >/dev/null 2>&1; do
  sleep 2
done

if ! ollama show "${MODEL}" >/dev/null 2>&1; then
  echo "Pulling Ollama model ${MODEL}..."
  ollama pull "${MODEL}"
else
  echo "Ollama model ${MODEL} already present."
fi

wait "${OLLAMA_PID}"
