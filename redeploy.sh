#!/bin/sh

curl --header "X-Webhook-Token: $REDEPLOY_TOKEN" http://213.183.51.37:9000/hooks/redeploy-webhook
