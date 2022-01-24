package templates

const TplStrHome = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Scenes of Shakespeare</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
	<link href="https://fonts.googleapis.com/css?family=Allura" rel="stylesheet">
  </head>
  <body>
    <div class="container">
	<div class="row" style="text-align: center">
		<h1 style="padding-top: 2em; font-family: Allura, cursive; font-size: 80px">Incidents Search</h1>
	</div>
	<div class="row" style="padding-top: 1.5em"><div class="col-sm-4 col-sm-offset-4">
		<form action="/" method="GET"><input class="form-control" autofocus name="q" maxlength=100 type="text"></form>
	</div></div>
	</div>
	{{range .Results}}
	<div class="row" style="padding: 1.5em 0; border-bottom: 1px solid #eee"><div class="col-sm-11 col-sm-offset-1">
		<div style="font-size: 1.2em">
			{{.IncidentId}}<br/>
			{{.DisplayName}}<br/>
			{{.Description}}<br/>
		</div>
	</div></div>
	{{end}}
    </div>
  </body>
</html>
`

const TplStrResults = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Incidents Search</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
	<link href="https://fonts.googleapis.com/css?family=Allura" rel="stylesheet">
  </head>
  <body>
    <div class="container-fluid">
	<div class="row" style="text-align: center; background-color: #f1f1f1; padding: 1em 0">
		<div class="col-sm-1" style="font-family: Allura, cursive; font-size: 30px">
		<a href="/" style="text-decoration: none; color: #000">Incidents</a>
		</div>
		<div class="col-sm-4">
			<form action="/" method="GET"><input class="form-control" autofocus name="q" maxlength=100 type="text" value="{{.Query}}"></form>
		</div>
	</div>
	{{range .Results}}
	<div class="row" style="padding: 1.5em 0; border-bottom: 1px solid #eee"><div class="col-sm-11 col-sm-offset-1">
		<div style="font-size: 1.2em">
			{{.IncidentId}}<br/>
			{{.DisplayName}}<br/>
			{{.Description}}<br/>
		</div>
	</div></div>
	{{end}}
	<div class="row" style="padding: 1.5em; background-color: #f1f1f1">
		Text and database from <a href="http://opensourceshakespeare.org/">http://opensourceshakespeare.org/</a>
    </div>
    </div>
  </body>
</html>
`