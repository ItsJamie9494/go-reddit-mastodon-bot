FROM scratch
COPY capybot /usr/bin/capybot
ENTRYPOINT [ "/usr/bin/capybot" ]