package htmxswap

import "github.com/gin-gonic/gin"

func AboutPage(c *gin.Context) {
	c.Header("Context-Type", "text/html")
	c.String(200, `
		<div class="container mt-5">
			<h2>About My Library</h2>
			<p class="mt-3">
				<strong>My Library</strong> is a web-based application designed to help readers discover, manage, and revisit books in a simple, organized way. Built around the <a href="https://openlibrary.org/developers/api" target="_blank">Open Library API</a>, it connects users with a vast collection of book and author data from around the world.
			</p>

			<h4 class="mt-4">🌟 What You Can Do</h4>
			<ul>
				<li><strong>Search & Explore:</strong> Look up any book or author using Open Library’s open database. Find covers, descriptions, publication details, and more.</li>
				<li><strong>Personal Library:</strong> Save books into your own categorized collection. Whether you're actively reading, planning your next read, or reminiscing a favorite, it's all in one place.</li>
				<li><strong>Track Progress:</strong> Organize your titles under:
					<ul>
						<li><strong>Reading</strong> – Books you're currently reading</li>
						<li><strong>Planning</strong> – Books you're planning to read</li>
						<li><strong>Finished</strong> – Books you've completed</li>
						<li><strong>Favorites</strong> – Books you especially loved</li>
					</ul>
				</li>
				<li><strong>Simple UI, Smart Updates:</strong> Uses <a href="https://htmx.org/" target="_blank">HTMX</a> to load content dynamically without full-page refreshes. It's fast and user-friendly.</li>
			</ul>

			<h4 class="mt-4">🛠️ Tech Stack</h4>
			<ul>
				<li><strong>Frontend:</strong> HTML, Bootstrap 5, HTMX</li>
				<li><strong>Backend:</strong> Golang, SQLite </li>
				<li><strong>Data Source:</strong> <a href="https://openlibrary.org/developers/api" target="_blank">Open Library API</a></li>
			</ul>

			<h4 class="mt-4">🎯 Project Goals</h4>
			<p>
				My Library was built with a passion for reading and the belief that managing books should be effortless and enjoyable. It aims to be a personal reading hub — not just a tracker, but a space that grows with you as a reader.
			</p>

			<h4 class="mt-4">🚀 What's Next?</h4>
			<ul>
				<li>User reviews & notes for books</li>
				<li>Custom reading challenges</li>
				<li>The ability to track where you are in the book/li>
				<li>Social features to share and follow other libraries</li>
			</ul>

			<h4 class="mt-4">🤝 Open Source & Contributions</h4>
			<p>
				This project is open source and under active development. If you’d like to contribute, suggest features, or report bugs, please visit the GitHub repository:
				<br>
				<a href="https://github.com/ctrl-MizteRy/Online-Library" target="_blank">https://github.com/ctrl-MizteRy/Online-Library</a>
			</p>

			<h5 class="mt-4">Thanks for stopping by, and happy reading! 📖</h5>
		</div>
			
        `)
}
