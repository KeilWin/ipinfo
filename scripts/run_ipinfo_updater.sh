#!/bin/bash

PROJECT_NAME="ipinfo_updater"

cp configs/$PROJECT_NAME.env bin/configs/$PROJECT_NAME.env

go run ./cmd/$PROJECT_NAME