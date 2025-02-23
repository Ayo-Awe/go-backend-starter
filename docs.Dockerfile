FROM redocly/redoc
EXPOSE 80
ENV SPEC_URL=openapi.yaml
COPY ./openapi.yaml /usr/share/nginx/html/openapi.yaml