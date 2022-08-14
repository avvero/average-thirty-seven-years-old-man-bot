#!/bin/sh

curl --header "X-Webhook-Token: $REDEPLOY_TOKEN" https://webhook.site/295f3d58-6f01-4b67-822f-fd69aa0d2e82
curl --header "X-Webhook-Token: $REDEPLOY_TOKEN" http://213.183.51.37:9000/hooks/redeploy-webhook