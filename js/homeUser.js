	
	var apiUrl = '/users/me/';

	$(document).ready(function(){
		
		// Les centres d'interets de l'utilisateur
		//Recuperer les elements
		$.ajax({
			url:	apiUrl + 'interests/',
			type:	'GET',
			success: function(data) {
				if (data) {
					$.each(data, function(index, value){
						$('#interests').append('<li class="interest"><a class="delete" data-id="'+ value.Id+'">x</a>  ' + value.Label + '</li>');
					})
				} else {
					console.log("ne s'interesse a rien pour le moment")
				}
			},
			error: function(xhr,status,error) {
				console.log(error);
			}
		});


		//Ajouter un Element d'interet a la liste
		$('#addinterest').on('submit', function(e){
			e.preventDefault();
			var content = $('#addinterest input[type="text"]').val();
			if (content.length > 0) {
				$.ajax({
					url:			apiUrl + 'interests/',
					type:			'POST',
					data:			JSON.stringify({Label:content}),
					contentType:	'application/json',
					dataType:		'json',
					success:		function(data) {
						if (data.Status == 'ok') {
							$('#addinterest input[type="text"]').val('');
						} else if (data.Status) {
							$('#interests').append('<li class="interest"><a class="delete" data-id="'+ data.Status+'">x</a>  ' + content + '</li>');
							$('#addinterest input[type="text"]').val('');
						} else {
							alert(data);
						}
					}
				})
			}
		});

		//Supprimer un element d'interet a la liste
		$('body').on('click', '.delete', function(e){
			e.preventDefault();
			var self = $(this);
			var id = self.attr('data-id');
			$.ajax({
				url:		apiUrl + 'interests/' + id,
				type:		'DELETE',
				success:	function(data) {
					if (data.Status == 'ok') {
						self.parent('.interest').remove();
					} else {
						alert(data);
					}
				}
			})
		});


		//Le num d'utilisateur du captaine
		//Recuperer l'age
		$.ajax({
			url:	apiUrl + 'username/',
			type:	'GET',
			success: function(data) {
				if (data) {
					console.log(data.Status)
					$('#changeusername input[type="text"]').val(data.Status)
				} else {
					console.log("What");
				}
			},
			error: function(xhr,status,error) {
				console.log(error);
			}
		});

		//Ajouter un Element d'interet a la liste
		$('#changeusername').on('submit', function(e){
			e.preventDefault();
			var content = $('#changeusername input[type="text"]').val();
			$.ajax({
				url:			apiUrl + 'username/',
				type:			'PUT',
				data:			JSON.stringify({Date:content}),
				contentType:	'application/json',
				dataType:		'json',
				success:		function(data) {
					if (data.Status == 'ok') {
						console.log('updated')
					} else {
						console.log("something wrong happen about your username");
					}
				},
				error: function(xhr,status,error) {
					console.log(error);
				}
			})
		});

		//Le prenom du captaine
		//Recuperer l'age
		$.ajax({
			url:	apiUrl + 'firstname/',
			type:	'GET',
			success: function(data) {
				if (data) {
					console.log(data.Status)
					$('#changeuFirstname input[type="text"]').val(data.Status)
				} else {
					console.log("What");
				}
			},
			error: function(xhr,status,error) {
				console.log(error);
			}
		});

		//Ajouter un Element d'interet a la liste
		$('#changeuFirstname').on('submit', function(e){
			e.preventDefault();
			var content = $('#changeuFirstname input[type="text"]').val();
			$.ajax({
				url:			apiUrl + 'firstname/',
				type:			'PUT',
				data:			JSON.stringify({Date:content}),
				contentType:	'application/json',
				dataType:		'json',
				success:		function(data) {
					if (data.Status == 'ok') {
						$('#userfirstname').empty();
						$('#userfirstname').append(content);
						console.log('updated')
					} else {
						console.log("something wrong happen about your username");
					}
				},
				error: function(xhr,status,error) {
					console.log(error);
				}
			})
		});

		//Le Nom de Famille du captaine
		//Recuperer le nom
		$.ajax({
			url:	apiUrl + 'lastname/',
			type:	'GET',
			success: function(data) {
				if (data) {
					console.log(data.Status)
					$('#changeuLastname input[type="text"]').val(data.Status)
				} else {
					console.log("What");
				}
			},
			error: function(xhr,status,error) {
				console.log(error);
			}
		});

		//Ajouter un Element d'interet a la liste
		$('#changeuLastname').on('submit', function(e){
			e.preventDefault();
			var content = $('#changeuLastname input[type="text"]').val();
			$.ajax({
				url:			apiUrl + 'lastname/',
				type:			'PUT',
				data:			JSON.stringify({Date:content}),
				contentType:	'application/json',
				dataType:		'json',
				success:		function(data) {
					if (data.Status == 'ok') {
						$('#userlastname').empty();
						$('#userlastname').append(content);
						console.log('updated')
					} else {
						console.log("something wrong happen about your username");
					}
				},
				error: function(xhr,status,error) {
					console.log(error);
				}
			})
		});

		//L'age du captaine
		//Recuperer l'age
		$.ajax({
			url:	apiUrl + 'age/',
			type:	'GET',
			success: function(data) {
				if (data) {
					console.log(data.Status)
					$('#changebirthdate input[type="date"]').val(data.Status)
				} else {
					console.log("quel age ??");
				}
			},
			error: function(xhr,status,error) {
				console.log(error);
			}
		});

		//Ajouter un Element d'interet a la liste
		$('#changebirthdate').on('submit', function(e){
			e.preventDefault();
			var content = $('#changebirthdate input[type="date"]').val();
			$.ajax({
				url:			apiUrl + 'age/',
				type:			'PUT',
				data:			JSON.stringify({Date:content}),
				contentType:	'application/json',
				dataType:		'json',
				success:		function(data) {
					if (data.Status == 'ok') {
						console.log('updated')
					} else {
						console.log("something wrong happen about your age");
					}
				},
				error: function(xhr,status,error) {
					console.log(error);
				}
			})
		});

		//Trouver des utilisateurs
		$.ajax({
			url:	apiUrl + 'matches/',
			type:	'GET',
			success: function(data) {
				if (data) {
					$.each(data, function(index, value){
						console.log(value);
						$('#matches').append('<li class="matche"><a href="users/' + value.Id + '">' + value.UserName + ' | ' + value.Bod + ' ans</a></li>');
					})
				} else {
					console.log("ne s'interesse a rien pour le moment")
				}
			},
			error: function(xhr,status,error) {
				console.log(error);
			}
		});


	});