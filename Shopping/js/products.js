class Products {
	constructor() {
		this.apiUrl = "http://localhost:8080/"; // Pointing to Go server
	}

	getNewProducts(limit) {
		$.ajax({
			type: "GET",
			url: this.apiUrl + "products?limit=" + limit + "&sort=desc", // API to fetch products from Go
			success: function (data) {
				$(data).each(function (index, product) {
					$(".products").append(
						'<div class="col-md-3"><div class="product"><a href="/product.html?productid=' +
							encodeURIComponent(product.id) +
							'">' +
							'<div class="image"><img src="' +
							product.image_url +
							'" class="img-fluid"></div><div class="info"><div class="title">' +
							product.name +
							"<br>$" +
							product.price +
							"</div></div></a></div></div>"
					);
				});
			},
		});
	}

	getSingleProduct(id) {
		$.ajax({
			type: "GET",
			url: this.apiUrl + "products/" + id, // API to fetch product details from Go
			success: function (data) {
				$(".breadcrumb").html(
					'<a href="/">Home</a><span class="sep">></span>' + data.name
				);
				$(".product_image").html('<img src="' + data.image_url + '" class="img-fluid">');
				$(".product_title").html(data.name);
				$(".product_price").html("$" + data.price.toFixed(2));
				$(".product_description").html("<p>" + data.description + "</p>");
			},
		});
	}
}
