#!/bin/bash
for filename in *.mjml; do
    mjml $filename -o ${filename%.mjml}.html
done