package main

import (
	"log"
	"net/http"
	"path/filepath"
	"flag"
	"os"
)

// Создаем структуру `application` для хранения зависимостей всего веб-приложения.
// Пока, что мы добавим поля только для двух логгеров, но
// мы будем расширять данную структуру по мере усложнения приложения.
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}
func main(){
	// Создаем новый флаг командной строки, значение по умолчанию: ":4000".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг. 
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера / HTTP")

	// Мы вызываем функцию flag.Parse() для извлечения флага из командной строки.
	// Она считывает значение флага из командной строки и присваивает его содержимое
	// переменной. Вам нужно вызвать ее *до* использования переменной addr
	// иначе она всегда будет содержать значение по умолчанию ":4000". 
	// Если есть ошибки во время извлечения данных - приложение будет остановлено.
	flag.Parse()

	// Используйте log.New() для создания логгера для записи информационных сообщений. Для этого нужно
	// три параметра: место назначения для записи логов (os.Stdout), строка
	// с префиксом сообщения (INFO или ERROR) и флаги, указывающие, какая
	// дополнительная информация будет добавлена. Обратите внимание, что флаги
	// соединяются с помощью оператора OR |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Создаем логгер для записи сообщений об ошибках таким же образом, но используем stderr как
	// место для записи и используем флаг log.Lshortfile для включения в лог
	// названия файла и номера строки где обнаружилась ошибка.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Инициализируем новую структуру с зависимостями приложения.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Инициализируем новую структуру http.Server. Мы устанавливаем поля Addr и Handler, так
	// что сервер использует тот же сетевой адрес и маршруты, что и раньше, и назначаем
	// поле ErrorLog, чтобы сервер использовал наш логгер
	// при возникновении проблем.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),	// Вызов нового метода app.routes()
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	// Вызываем метод ListenAndServe() от нашей новой структуры http.Server
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
type neuteredFileSystem struct {
	fs http.FileSystem
}
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
	
}