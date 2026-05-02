"""
Ghost Hunter AI Training - Numpy Only Version
Simple neural network trained on your gameplay data.
"""

import json
import numpy as np
import os

# Load data
print("Loading training data...")
with open("training_data.json", "r") as f:
    data = json.load(f)

print(f"Loaded {len(data)} samples")

# Parse state to features
def parse_state(s):
    features = [
        s["player_x"] / 32.0,
        s["player_y"] / 32.0,
        s["player_angle"] / (2 * np.pi),
        s["health"] / 100.0,
        s["ammo"] / 50.0,
        s["weapon"] / 2.0,
        s["enemy_count"] / 20.0,
        s["wave"] / 5.0,
        s["current_map"] / 4.0,
        s["portal_dist"] / 20.0,
        s["ammo_pickup_dist"] if s["has_ammo_pickup"] else 99.0,
        s["health_pickup_dist"] if s["has_health_pickup"] else 99.0,
    ]
    # Enemy features (up to 5)
    for i in range(5):
        if i < len(s["enemy_distances"]):
            features.extend([s["enemy_distances"][i]/15.0, s["enemy_angles"][i]/np.pi])
        else:
            features.extend([1.0, 0.0])
    return np.array(features, dtype=np.float32)

def parse_action(a):
    return np.array([
        float(a["move_forward"]),
        float(a["move_backward"]),
        float(a["turn_left"]),
        float(a["turn_right"]),
        float(a["shoot"])
    ], dtype=np.float32)

# Parse all data
X = np.array([parse_state(d["state"]) for d in data])
y = np.array([parse_action(d["action"]) for d in data])

print(f"Features: {X.shape}, Actions: {y.shape}")

# Normalize
X_mean = X.mean(axis=0)
X_std = X.std(axis=0) + 1e-8
X = (X - X_mean) / X_std

# Simple Neural Network class
class SimpleNN:
    def __init__(self, input_size, output_size, hidden=64):
        # Xavier initialization
        self.W1 = np.random.randn(input_size, hidden) * np.sqrt(2.0/input_size)
        self.b1 = np.zeros(hidden)
        self.W2 = np.random.randn(hidden, hidden) * np.sqrt(2.0/hidden)
        self.b2 = np.zeros(hidden)
        self.W3 = np.random.randn(hidden, output_size) * np.sqrt(2.0/hidden)
        self.b3 = np.zeros(output_size)
    
    def relu(self, x):
        return np.maximum(0, x)
    
    def sigmoid(self, x):
        return 1 / (1 + np.exp(-np.clip(x, -500, 500)))
    
    def forward(self, x):
        self.z1 = x @ self.W1 + self.b1
        self.a1 = self.relu(self.z1)
        self.z2 = self.a1 @ self.W2 + self.b2
        self.a2 = self.relu(self.z2)
        self.z3 = self.a2 @ self.W3 + self.b3
        return self.sigmoid(self.z3)
    
    def train(self, X, y, epochs=50, lr=0.01, batch_size=64):
        n = len(X)
        for epoch in range(epochs):
            # Shuffle
            idx = np.random.permutation(n)
            total_loss = 0
            
            for i in range(0, n, batch_size):
                batch_idx = idx[i:i+batch_size]
                X_batch = X[batch_idx]
                y_batch = y[batch_idx]
                
                # Forward
                out = self.forward(X_batch)
                
                # Loss (MSE)
                loss = np.mean((out - y_batch) ** 2)
                total_loss += loss
                
                # Backward (gradient descent)
                delta3 = (out - y_batch) * out * (1 - out)  # sigmoid derivative
                dW3 = self.a2.T @ delta3 / batch_size
                db3 = delta3.mean(axis=0)
                
                delta2 = (delta3 @ self.W3.T) * (self.a2 > 0)
                dW2 = self.a1.T @ delta2 / batch_size
                db2 = delta2.mean(axis=0)
                
                delta1 = (delta2 @ self.W2.T) * (self.a1 > 0)
                dW1 = X_batch.T @ delta1 / batch_size
                db1 = delta1.mean(axis=0)
                
                # Update
                self.W3 -= lr * dW3
                self.b3 -= lr * db3
                self.W2 -= lr * dW2
                self.b2 -= lr * db2
                self.W1 -= lr * dW1
                self.b1 -= lr * db1
            
            if (epoch + 1) % 10 == 0:
                print(f"Epoch {epoch+1}/{epochs}, Loss: {total_loss/(n/batch_size):.4f}")
    
    def save(self, filename):
        np.savez(filename, 
                 W1=self.W1, b1=self.b1,
                 W2=self.W2, b2=self.b2,
                 W3=self.W3, b3=self.b3,
                 mean=X_mean, std=X_std)

# Train
print("\nTraining...")
model = SimpleNN(X.shape[1], y.shape[1], hidden=64)
model.train(X, y, epochs=50, lr=0.01)

# Save
model.save("ghost_hunter_ai.npz")
print("\nModel saved to ghost_hunter_ai.npz")

# Export as JSON for Go - flatten arrays
print("Exporting as JSON...")
model_data = {
    "input_size": int(X.shape[1]),
    "output_size": int(y.shape[1]),
    "W1": model.W1.tolist(),  # Already list of lists
    "b1": model.b1.tolist(),
    "W2": model.W2.tolist(),
    "b2": model.b2.tolist(),
    "W3": model.W3.tolist(),
    "b3": model.b3.tolist(),
    "mean": X_mean.tolist(),
    "std": X_std.tolist()
}

# Pretty print for debugging
with open("model_weights.json", "w") as f:
    json.dump(model_data, f, indent=2)

print("Done! model_weights.json created")