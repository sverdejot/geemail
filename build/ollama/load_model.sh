#!/bin/bash

ollama serve & 

echo "Ollama is ready, creating the model..."

ollama create geemail -f /model_files/Modelfile
ollama run geemail
