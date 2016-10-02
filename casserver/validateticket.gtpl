<html>
<head>
<title></title>
</head>
<body>
Ticket: {{.Ticket}}
<form action="{{.CallbackURL}}" method="post">
    AuthenticationSuccess: CasUser: <input type="text" name="casuser" /><br />
    AuthenticationFailure: <input type="text" name="authFailure" /><br />
    FailureReason: <input type="text" name="failureReason" /><br />
    <input type="submit" value="Send">
</form>
</body>
</html>