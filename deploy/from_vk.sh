#!/bin/bash

yc serverless function version create \
  --function-name from-vk \
  --runtime golang123 \
  --entrypoint index.Handler \
  --memory 128m \
  --execution-timeout 15s \
  --environment TELEGRAM_CHAT_ID=-1003003844767,VK_GROUP_ID=232485584 \
  --source-path povtorusca-bot \
  --service-account-id aje1qaecproapgu5anqm \
  --secret id=e6qomlnufdidl6n2gumf,key=confirmation_code,environment-variable=VK_CONFIRMATION_CODE \
  --secret id=e6qomlnufdidl6n2gumf,key=secret,environment-variable=VK_SECRET \
  --secret id=e6q6lneuvm8e6tnhtvos,key=token,environment-variable=TELEGRAM_BOT_TOKEN \
  --concurrency 5
