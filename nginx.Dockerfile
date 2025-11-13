FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=frontend /app/dist /usr/share/nginx/html
