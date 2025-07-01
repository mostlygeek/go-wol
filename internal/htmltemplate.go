package internal

// HTML template for the web interface
var TEMPLATE = `
<!DOCTYPE html>
<html>
<head>
    <title>Wake-on-LAN</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 100%;
            padding: 20px;
            box-sizing: border-box;
            display: flex;
            flex-direction: column;
            align-items: center;
            min-height: 100vh;
            margin: 0;
        }
        .container {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            flex-grow: 1;
            width: 100%;
            max-width: 400px;
            margin: 0 auto;
        }
        #address {
            width: 100%;
            padding: 12px 20px;
            margin: 8px 0;
            box-sizing: border-box;
            border: 2px solid #ccc;
            border-radius: 4px;
            font-size: 16px;
        }


        .wake-button {
            background-color: #00dd00;
            color: white;
            border: none;
            border-radius: 50%;
            padding: 40px;
            font-size: 24px;
            cursor: pointer;
            width: 150px;
            height: 150px;
            margin: 20px 0;
            box-shadow: 0 0 0 8px #333,
                        0 0 0 12px #222,
                        inset 0 -10px 15px rgba(0,0,0,0.4),
                        inset 0 10px 15px rgba(255,255,255,0.3);
            transition: all 0.1s ease;
            position: relative;
            font-weight: bold;
            text-transform: uppercase;
            letter-spacing: 1px;
            text-shadow: 0 1px 3px rgba(0,0,0,0.5);
        }

        .wake-button:hover {
            background-color: #00ff00;
            transform: scale(1.05);
            box-shadow: 0 0 0 8px #333,
                        0 0 0 12px #222,
                        0 0 20px rgba(0,255,0,0.4),
                        inset 0 -10px 15px rgba(0,0,0,0.4),
                        inset 0 10px 15px rgba(255,255,255,0.3);
        }

        .wake-button:active {
            transform: scale(0.97);
            background-color: #00aa00;
            box-shadow: 0 0 0 8px #333,
                        0 0 0 12px #222,
                        inset 0 -5px 10px rgba(0,0,0,0.6),
                        inset 0 5px 10px rgba(255,255,255,0.2);
        }

        /* Add an arcade button shine effect */
        .wake-button::before {
            content: '';
            position: absolute;
            top: 15%;
            left: 15%;
            width: 30%;
            height: 30%;
            background: rgba(255,255,255,0.3);
            border-radius: 50%;
            pointer-events: none;
        }


        h1 {
            text-align: center;
            color: #333;
        }

         .message {
            margin: 15px 0;
            padding: 10px;
            border-radius: 4px;
            width: 100%;
            text-align: center;
        }
        .success {
            background-color: #dff0d8;
            color: #3c763d;
            display: block;
        }
        .error {
            background-color: #f2dede;
            color: #a94442;
            display: block;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Wake-on-LAN</h1>
        <input type="text" id="address" placeholder="MAC Address to send to" value="88:88:88:88:87:88">
        <button class="wake-button" onclick="sendWakePacket()">Wake Up</button>
        <div id="message" class="message">Press Button</div>
    </div>

    <script>
        // Load saved address from localStorage
        document.addEventListener('DOMContentLoaded', function() {
            const savedAddress = localStorage.getItem('wolAddress');
            if (savedAddress) {
                document.getElementById('address').value = savedAddress;
            }
        });

        function showMessage(text, isError) {
            const messageEl = document.getElementById('message');
            messageEl.textContent = text;
            messageEl.className = isError ? 'message error' : 'message success';
        }

        function sendWakePacket() {
            const address = document.getElementById('address').value;
            if (!address) {
                showMessage('Please enter a MAC address', true);
                return;
            }

            // Save address to localStorage
            localStorage.setItem('wolAddress', address);

            fetch('/wakeup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ address: address }),
            })
            .then(response => response.json())
            .then(data => {
                showMessage(data.message || 'Magic packet sent successfully!', false);
            })
            .catch(error => {
                showMessage('Error sending magic packet: ' + error.message, true);
            });
        }
    </script>
</body>
</html>`
