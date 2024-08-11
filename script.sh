#!/bin/bash

PROGRAM_PATH="./s7parser"

TARGET_TIMES=("08:00" "13:00" "16:00")

while true; do
    CURRENT_TIME=$(date +"%H:%M")

    for TARGET_TIME in "${TARGET_TIMES[@]}"; do
        if [ "$CURRENT_TIME" == "$TARGET_TIME" ]; then
            "$PROGRAM_PATH"
            TARGET_TIMES=("${TARGET_TIMES[@]/$TARGET_TIME}")
        fi
    done

    if [ ${#TARGET_TIMES[@]} -eq 0 ]; then
        break
    fi

    sleep 10
done
