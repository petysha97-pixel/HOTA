package handlers_http

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// регистрация пользователей через POST запросы
func NewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(405)
		return
	}

	if r.Method == "POST" {
		var user models.User

		body, err := io.ReadAll(r.Body)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Не верные данные 1 ")
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Не верные данные 2")
			return
		}
	

		errs := validation(user)
		if len(errs) > 0 {
			status := models.Registration_fail{
				StatusGlobal: "Регистрация не успешна",
				Error:        errs,
			}
			data, _ := json.MarshalIndent(&status,
				" ",
				" ")

			w.Write(data)
			return
		}
		if len(errs) == 0 {
			status := models.Registration_fail{
				StatusGlobal: "Регистрация успешна",
			}
			data, _ := json.MarshalIndent(&status,
				" ",
				" ")

			w.Write(data)
				
		}

		user = repositories.AppendUser(user)
		fmt.Println(user)
	}
}

// валидация полей
func validation(user models.User) []string {

	var result []string

	if user.Nickname == "" {
		err := "Пустой никнейм"
		result = append(result, err)
	}
	if user.Role == "" {
		err := "Пустая роль"
		result = append(result, err)
	}
	if len(user.Stack) == 0 {
		err := "Пустой стек"
		result = append(result, err)
	}
	if user.GitHub == "" {
		err := "Пустой гитхаб"
		result = append(result, err)
	}
	if user.Telegram == "" {
		err := "Пустой ТГ"
		result = append(result, err)
	}
	if user.Status == "" {
		err := "Пустой статус"
		result = append(result, err)
	}

	if len(result) > 0 {
		return result
	}

	return nil
}


func Home(w http.ResponseWriter, r *http.Request) {

	html := `<!DOCTYPE html>
<html lang="ru">

<head>
<meta charset="UTF-8">
<title>НОТА</title>

<style>

* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
	font-family: Segoe UI, sans-serif;
}

body {
	height: 100vh;

	background:
	linear-gradient(
		135deg,
		#050505,
		#0b1611,
		#10261d
	);

	display: flex;
	justify-content: center;
	align-items: center;
}

.card {

	width: 550px;

	padding: 40px;

	background: rgba(15,42,32,0.95);

	border: 1px solid #6e5d26;

	border-radius: 20px;

	box-shadow:
		0 0 30px rgba(214,177,61,0.15);

}

.logo {

	text-align: center;

	font-size: 36px;

	font-weight: bold;

	letter-spacing: 12px;

	color: #d6b141;

	margin-bottom: 10px;
}

.subtitle {

	text-align: center;

	color: #7f8c84;

	font-size: 13px;

	letter-spacing: 3px;

	margin-bottom: 35px;
}

label {

	display: block;

	color: #d6d6d6;

	margin-bottom: 8px;

	font-size: 14px;
}

input {

	width: 100%;

	padding: 14px;

	margin-bottom: 18px;

	border-radius: 10px;

	border: 1px solid #30463c;

	background: #14251d;

	color: white;

	outline: none;
}

input:focus {

	border-color: #d6b13d;
}

button {

	width: 100%;

	padding: 15px;

	border: none;

	border-radius: 10px;

	background: #d6b13d;

	color: black;

	font-size: 15px;

	font-weight: bold;

	cursor: pointer;

	transition: 0.3s;
}

button:hover {

	transform: translateY(-2px);

	background: #e8c24d;
}

.footer {

	text-align: center;

	margin-top: 20px;

	color: #6d7a74;

	font-size: 12px;
}

#result {

	margin-top: 20px;

	color: #ffffff;

	white-space: pre-wrap;
}

</style>

</head>

<body>

<div class="card">

	<div class="logo">
		НОТА
	</div>

	<div class="subtitle">
		ПОДКЛЮЧЕНИЕ УЗЛА К СЕТИ
	</div>

	<form id="userForm">

		<label>Ник</label>
		<input id="nickname">

		<label>Роль</label>
		<input id="role">

		<label>Стек</label>
		<input id="stack" placeholder="go,docker,postgres">

		<label>GitHub</label>
		<input id="github">

		<label>Telegram</label>
		<input id="telegram">

		<button type="submit">
			ПОДКЛЮЧИТЬСЯ
		</button>

	</form>

	<div id="result"></div>

	<div class="footer">
		HOTA NETWORK © 2026
	</div>

</div>

<script>

document
.getElementById("userForm")
.addEventListener("submit", async function(e){

	e.preventDefault();

	const data = {
		nickname: document.getElementById("nickname").value,
		role: document.getElementById("role").value,
		stack: document.getElementById("stack").value.split(","),
		github: document.getElementById("github").value,
		telegram: document.getElementById("telegram").value,
		status: "online"
	};

	const response = await fetch("/user",{
		method:"POST",
		headers:{
			"Content-Type":"application/json"
		},
		body:JSON.stringify(data)
	});

	const result = await response.json();

	if(result.error){

		let text = "";

		result.error.forEach(err=>{
			text += "❌ " + err + "\n";
		});

		document.getElementById("result").innerText = text;

	}else{

		window.location.href = "/profile";

	}

});

</script>

</body>
</html>`

	fmt.Fprint(w, html)
}


func Profile(w http.ResponseWriter, r *http.Request) {

	user, ok := repositories.LastUser()
	if !ok {
		http.Error(w, "Пользователь не найден", 404)
		return
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Профиль</title>
</head>

<body style="background:#0b1611;color:white;font-family:Segoe UI;">

<h1>%s</h1>

<p>Роль: %s</p>

<p>GitHub: %s</p>

<p>Telegram: %s</p>

<p>Статус: %s</p>

<p>Стек:</p>

<ul>
`, user.Nickname,
		user.Role,
		user.GitHub,
		user.Telegram,
		user.Status)

	for _, tech := range user.Stack {
		html += fmt.Sprintf("<li>%s</li>", tech)
	}

	html += `
</ul>

</body>
</html>
`

	fmt.Fprint(w, html)
}