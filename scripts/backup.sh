#!/bin/bash
POD=$(kubectl get pods -l app.kubernetes.io/name=divelog -o jsonpath="{.items[0].metadata.name}")
kubectl cp $POD:/app/data/divelog.db ./data/backups/divelog-$(date +%Y%m%d-%H%M%S).db
echo "Backup complete"