FROM centurylink/ca-certs
ADD main /
EXPOSE 8080
CMD ["/main"]