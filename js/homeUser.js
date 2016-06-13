	
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

		//L'age du captaine
		//Recuperer l'age
		$.ajax({
			url:	apiUrl + 'age/',
			type:	'GET',
			success: function(data) {
				if (data) {
					var d = String(data.Date).split();
					console.log(d[0]);
					$('#changebirthdate input[type="text"]').val(d[0])
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


	});