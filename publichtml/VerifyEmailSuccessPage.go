package publichtml

// VerifyEmailSuccessPage ...
var VerifyEmailSuccessPage = `
<!DOCTYPE html>
<html>
  <!-- This page is displayed whenever someone has successfully reset their password.
       Pro and Enterprise accounts may edit this page and tell Parse to use that custom
       version in their Parse app. See the App Settigns page for more information.
       This page will be called with the query param 'username'
   -->
  <head>
  <title>Email Verification</title>
  <style type='text/css'>
    h1 {
      color: #0067AB;
      display: block;
      font: inherit;
      font-family: 'Open Sans', 'Helvetica Neue', Helvetica;
      font-size: 30px;
      font-weight: 600;
      height: 30px;
      line-height: 30px;
      margin: 45px 0px 0px 45px;
      padding: 0px 8px 0px 8px;
    }
  </style>
  <body>
    <h1>Successfully verified your email!</h1>
  </body>
</html>
`
