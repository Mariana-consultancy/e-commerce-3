#Project Title: Sole Select

**Description**: Contributed and collaborated with project managers and developers on an e-commerce project called Sole Select - a sneaker trading website.

Backend: Go, Gin
Database: PostgreSQL, GORM
Deployment:

Backend Deployed on: Render
Tools & Project Management
Version Control: GitHub

Project Management:
Jira: We used Jira to assign and manage tasks with the product manager and other developers. We used Jira to plan a 4-week sprint to complete the e-commerce project and updated Jira daily after stand-ups.
Slack: We used Slack as a communication platform, by communicating in huddles and channels to share findings, progress and ideas.

Features:

Set up user authentication using JWT-based authentication.
Users can place orders, view their order history, and check order statuses.
CRUD APIs.
Key Features:

Salary and Expense Management: Users can submit their monthly salary, and expenses are validated in real-time to ensure they don't exceed the remaining salary. The app provides an overview of all expenses, including the total amount spent and the remaining balance.
CRUD Operations: Full support for creating, reading, updating, and deleting expenses using a RESTful API designed with the Gin framework, integrated with GORM for database interaction.
Error Handling and Validation: Error handling and data validation to ensure accurate and correct data entry, providing clear feedback for invalid inputs using status codes and responding with a JSON response.
My role:

Backend Development: Designed and implemented the backend API using the Gin web framework in Go, ensuring efficient routing and integration with the database for salary and expense management features.
Database Integration: Integrated GORM for handling database operations, including salary and expense CRUD functionality, ensuring smooth and reliable data storage and retrieval.
Challenges and Solutions:

Challenge 1: Ensuring that expenses do not exceed the remaining salary, and handling cases where users input incorrect data.

Solution 1: Implemented data validation using Gin’s bind JSON and custom error handling to provide clear error messages and maintain data consistency.

Challenge 2: Users must input their monthly salary before entering expenses to prevent errors.

Solution 2: Implemented a check for salary requirements before inputting expenses. An error will be displayed if expenses are entered first.

Results and Impact: Achieved a 95% satisfaction rate from 5 testing users.

Links: Live Demo GitHub Repository

Lessons Learned: Improved skills in backend optimisation and gained experience using RESTful API with GORM for database interaction. For example, retrieving data from the database using gorm first and find.

Lessons Learned: Improved skills in backend optimisation and gained experience using RESTful API with GORM for database interaction.

Full API Documentation
