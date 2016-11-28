package main

import "net/http"

func router() {
	instance, err := mgoSession()

	if err != nil {
		logger.Println("cant connect to database")
		// TODO: email panic
		return
	}

	db := instance.session
	ensureIndexes(db)

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/forgotpass", ForgotPasswordHandler)

	base := []Adapter{
		logRequest(logger),
		mongo(db),
	}

	common := []Adapter{
		logRequest(logger),
		mongo(db),
		authenticate(),
		currentUser(),
	}

	http.Handle("/update-conditions", chain(
		httpHandler(updateConditionHandler),
		common...,
	))

	http.Handle("/delete-condition", chain(
		httpHandler(deleteConditionHandler),
		common...,
	))

	http.Handle("/get-conditions", chain(
		httpHandler(getConditionsIndexHandler),
		common...,
	))

	http.Handle("/create-condition", chain(
		httpHandler(createConditionHandler),
		common...,
	))

	http.Handle("/update-calls", chain(
		httpHandler(updateCallHandler),
		common...,
	))

	http.Handle("/delete-calls", chain(
		httpHandler(deleteCallHandler),
		common...,
	))

	http.Handle("/get-calls", chain(
		httpHandler(getCallsIndexHandler),
		common...,
	))

	http.Handle("/create-call", chain(
		httpHandler(createCallHandler),
		common...,
	))

	http.Handle("/userprofile", chain(
		httpHandler(UserProfileHandler),
		logRequest(logger),
		mongo(db),
		authenticate(),
		currentUser(),
	))

	http.Handle("/requestpass", chain(
		httpHandler(RequestPasswordHandler),
		base...,
	))

	http.Handle("/cp", chain(
		httpHandler(ChangePassFormHandler),
		logRequest(logger),
		mongo(db),
	))

	http.Handle("/changepass", chain(
		httpHandler(ChangePasswordHandler),
		logRequest(logger),
		mongo(db),
	))

	// Handle POST/GET/DELETE/PUT to /users
	http.Handle("/users", chain(
		// http handler for every request to /users
		httpHandler(UsersHandler),

		// log request to stdout, like: GET /users
		logRequest(logger),

		// Copy db session, it will take care of close connection using defer
		mongo(db),

		// Authenticate local user
		authenticate(),

		// Load current user into memory
		currentUser(),
	))

	// Handle POST/GET/DELETE/PUT to /users
	http.Handle("/edit-user", chain(
		// http handler for every request to /users
		httpHandler(editUserHandler),

		// log request to stdout, like: GET /users
		logRequest(logger),

		// Copy db session, it will take care of close connection using defer
		mongo(db),
	))

	// Handle POST/GET/DELETE/PUT to /users
	http.Handle("/update-user", chain(
		// http handler for every request to /users
		httpHandler(UpdateUserHandler),

		// log request to stdout, like: GET /users
		logRequest(logger),

		// Copy db session, it will take care of close connection using defer
		mongo(db),
	))

	// Handle POST/GET /clients
	http.Handle("/clients", chain(
		// http handler for every request to /users
		httpHandler(ClientsHandler),

		// log request to stdout, like: GET /users
		logRequest(logger),

		// Copy db session, it will take care of close connection using defer
		mongo(db),

		// Authenticate local user
		authenticate(),

		// Load current user into memory
		currentUser(),
	))

	// Handle PUT Clients
	http.Handle("/edit-client", chain(
		httpHandler(EditClientHandler),
		common...,
	))

	// Handle DELETE Clients
	http.Handle("/delete-client", chain(
		httpHandler(DeleteClientHandler),
		common...,
	))

	// Handle POST/GET /clients
	http.Handle("/client-wizard", chain(
		// http handler for every request to /users
		httpHandler(ClientWizardHandler),

		// log request to stdout, like: GET /users
		logRequest(logger),

		// Copy db session, it will take care of close connection using defer
		mongo(db),
	))

	http.Handle("/update-client", chain(
		httpHandler(UpdateClientHandler),
		common...,
	))

	// Handle POST/GET /clients
	http.Handle("/scheduling-logic", chain(
		httpHandler(SchedulingLogicHandler),
		common...,
	))

	// Handle POST /schedule
	http.Handle("/schedule", chain(
		httpHandler(createScheduleHandler),
		common...,
	))

	// Handle POST /exception
	http.Handle("/exception", chain(
		httpHandler(createExceptionHandler),
		common...,
	))

	http.Handle("/delete-exception", chain(
		httpHandler(deleteExceptionHandler),
		common...,
	))

	http.Handle("/update-exception", chain(
		httpHandler(updateExceptionHandler),
		common...,
	))

	// Handle GET /confirm
	http.Handle("/confirm", chain(
		// http handler for every request to /users
		httpHandler(ConfirmHandler),

		// log request to stdout, like: GET /users
		logRequest(logger),

		// Copy db session, it will take care of close connection using defer
		mongo(db),
	))

	// Handle GET /confirm
	http.Handle("/set-pass", chain(
		// http handler for every request to /users
		httpHandler(PasswordHandler),

		// log request to stdout, like: GET /users
		logRequest(logger),

		// Copy db session, it will take care of close connection using defer
		mongo(db),
	))

	// Handle /login
	http.Handle("/login", chain(
		httpHandler(LoginHandler),
		mongo(db)),
	)

	// Serve assets: for now using go, on production we can use nginx
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
}
