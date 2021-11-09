package main

const listFilmsHTML = `
<Html>
	<Head>
		<title>Список фильмов</title>
	</Head>

	<Body>
		<h1>Films - {{.Count}}</h1>
		<table>
			<tr style='text-align: left'>
				<th>Name</th>
				<th>Year</th>
				<th>Genre</th>
				<th>Director</th>
			</tr>
			{{range .Films}}
			<tr>
				<td><a href='{{.Title | wrapReplaceAllSpace}}_{{.Year}}.html'>{{.Title}}</a></td>
				<td>{{.Year}}</td>
				<td>{{.Genre}}</td>
				<td>{{.Director}}</td>
			</tr>
			{{end}}
		</table>
	</Body>  
</Html>
`
