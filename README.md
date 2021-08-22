# Setting up your Azurre application

1. Navigate and sign into portal.azure.com > Azure Active Directory
2. Select "App Registrations" > "New Registration"
3. Create a Client Secret
4. Navigate to Authentication > Add a Platform > Web
    - For local testing, my redirect URI poitns to `http://localhost/auth/microsoft/callback`, which you should update to the real redirect endpoint after the app goes live
    - Select ID Tokens for the implicit grant and hybrid flows
    - Select "Accounts in this organizational directory only"
    - Enable the following mobile and desktop flows:
        - App collects plaintext password (Resource Owner Password Credential Flow) Learn more
        - No keyboard (Device Code Flow) Learn more
        - SSO for domain-joined Windows (Windows Integrated Auth Flow) Learn more
    - Save your configuration

# Running the application
Create a `.env` at the root directory with the same format as the sample `.sample-env`.
Run your application.
Nagivate to `http://localhost:8000/`, login with your Azure AD account, and see the token.