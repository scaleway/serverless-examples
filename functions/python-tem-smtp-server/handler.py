import os
import logging
import smtplib
from email.mime.text import MIMEText


def handle(event, context):

    # SMTP configuration, see https://www.scaleway.com/en/docs/managed-services/transactional-email/reference-content/smtp-configuration/ for more info

    host = "smtp.tem.scw.cloud"  # The domain name or IP address of the SMTP server. If you are using Scaleway TEM, the domain to use is smtp.tem.scw.cloud
    port = 465

    # Your Scaleway SMTP username is the Project ID of the Project in which the TEM domain was created
    login = os.environ.get("TEM_PROJECT_ID")
    # Your password is the secret key of the API key of the project used to manage your TEM domain
    password = os.environ.get("SECRET_KEY")

    # the email address that will appear as the sender
    sender_email = "jdoe@fake-domain.com"
    # the email address to send the email to
    receiver_email = "jdoe@fake-email.com"

    message = MIMEText(
        "You've successfully sent an email from Serverless Functions!", "plain"
    )
    message["Subject"] = "Congratulations"
    message["From"] = sender_email
    message["To"] = receiver_email

    try:
        with smtplib.SMTP_SSL(host, port) as server:
            server.login(login, password)
            server.sendmail(sender_email, receiver_email, message.as_string())
    except Exception as e:
        logging.error("Failed to send email: %s", e)


if __name__ == "__main__":
    from scaleway_functions_python import local

    local.serve_handler(handle)
