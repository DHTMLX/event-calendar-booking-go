FROM debian:bookworm-slim
WORKDIR /app
COPY ./event-calendar-booking /app/event-calendar-booking

EXPOSE 3000

CMD ["/app/event-calendar-booking"]
