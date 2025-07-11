const username = prompt("Enter your username:") || "Anonymous";
const socket = new WebSocket(`ws://localhost:8080/ws?username=${encodeURIComponent(username)}`);

const messageForm = document.getElementById("message-form");
const messageInput = document.getElementById("message");
const messagesContainer = document.getElementById("messages");

socket.onopen = () => {
  console.log("Connected to server");
};

socket.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  const div = document.createElement("div");
  div.textContent = `${msg.username}: ${msg.content}`;
  messagesContainer.appendChild(div);
};

socket.onerror = (err) => {
  console.error("WebSocket error:", err);
};

messageForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const content = messageInput.value.trim();
  if (content) {
    const message = {
      username: username,
      content: content
    };
    socket.send(JSON.stringify(message));
    messageInput.value = "";
  }
});