FROM registry.suse.com/bci/python:3.13

WORKDIR /app

COPY requirements.txt ./
RUN pip uninstall openai
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD [ "python3", "./openai-client-otel.py" ]
