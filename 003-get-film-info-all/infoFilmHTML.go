package main

const infoFilmHTML = `
<Html>
	<Head>
		<title>{{.Title}}</title>
	</Head>

	<Body>
		<h1>{{.Title}}</h1>
		<img src="../images/Poster_{{.Title | wrapReplaceAllSpace}}_{{.Year}}.jpg" alt="{{.Title}}">

		<h2>О фильме</h2>
		<table>
			<tr>
				<td>Год производства</td>
				<td>{{.Year}}</td>
			</tr>
			<tr>
				<td>Релиз</td>
				<td>{{.Released}}</td>
			</tr>
			<tr>
				<td>Жанр</td>
				<td>{{.Genre}}</td>
			</tr>
			<tr>
				<td>Режиссер</td>
				<td>{{.Director}}</td>
			</tr>
			<tr>
				<td>Cценарий</td>
				<td>{{.Writer}}</td>
			</tr>
			<tr>
				<td>Актеры</td>
				<td>{{.Actors}}</td>
			</tr>
			<tr>
				<td>Сборы в США</td>
				<td>{{.BoxOffice}}</td>
			</tr>
			<tr>
				<td>Страна</td>
				<td>{{.Country}}</td>
			</tr>
			<tr>
				<td>Награды</td>
				<td>{{.Awards}}</td>
			</tr>
			<tr>
				<td>Язык оригинала</td>
				<td>{{.Language}}</td>
			</tr>
			<tr>
				<td>Тип</td>
				<td>{{.Type}}</td>
			</tr>
			<tr>
				<td>Время</td>
				<td>{{.Runtime}}</td>
			</tr>
			<tr>
				<td>Рейтинг MPAA</td>
				<td>{{.Rated}}</td>
			</tr>
		</table>
		<br>

		<h3>Описание:</h3>
		{{.Plot}}<br><br>

		<h3>Рейтинги</h3>
		<table>
			<tr style='text-align: left'>
				<th>Source</th>
				<th>Value</th>
			</tr>
			{{range .Ratings}}
			<tr>
				<td>{{.Source}}</td>
				<td>{{.Value}}</td>
			</tr>
			{{end}}
		</table>
		<br>

		<a href="https://www.imdb.com/title/{{.ImdbID}}/">Internet Movie Database</a>
		<br>
		<a href="{{.Poster}}">Ссылки на постеры</a>
	</Body>  
</Html> 
`
