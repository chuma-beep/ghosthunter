"""
Ghost Hunter AI Training Script
Trains a neural network to play the game using imitation learning.
"""

import json
import numpy as np
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader
from collections import Counter
import os

# Check for training data
if not os.path.exists("training_data.json"):
    print("No training data found!")
    print("Run the game and press D to start collecting data, D again to stop.")
    print("Play the game for a while to generate training samples.")
    exit(1)

# Load training data
print("Loading training data...")
with open("training_data.json", "r") as f:
    data = json.load(f)

print(f"Loaded {len(data)} samples")

# Convert to numpy arrays
def parse_state(sample):
    """Extract features from game state"""
    s = sample["state"]
    a = sample["action"]
    
    features = [
        s["player_x"] / 32.0,  # Normalize to 0-1 (map size)
        s["player_y"] / 32.0,
        s["player_angle"] / (2 * np.pi),  # Normalize angle
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
    
    # Add enemy features (up to 5 enemies)
    for i in range(5):
        if i < len(s["enemy_distances"]):
            features.extend([
                s["enemy_distances"][i] / 15.0,
                s["enemy_angles"][i] / np.pi,
            ])
        else:
            features.extend([1.0, 0.0])  # No enemy
    
    return np.array(features, dtype=np.float32)

def parse_action(sample):
    """Convert action to one-hot or multi-label"""
    a = sample["action"]
    
    # Encode as multi-label (can have multiple actions)
    actions = [
        float(a["move_forward"]),
        float(a["move_backward"]),
        float(a["turn_left"]),
        float(a["turn_right"]),
        float(a["shoot"]),
    ]
    
    return np.array(actions, dtype=np.float32)

# Parse all data
X = np.array([parse_state(d) for d in data])
y = np.array([parse_action(d) for d in data])

print(f"Feature shape: {X.shape}")
print(f"Action shape: {y.shape}")

# Analyze action distribution
action_names = ["Forward", "Backward", "Left", "Right", "Shoot"]
for i, name in enumerate(action_names):
    count = int(y[:, i].sum())
    print(f"{name}: {count} ({100*count/len(y):.1f}%)")

# Normalize features
X_mean = X.mean(axis=0)
X_std = X.std(axis=0) + 1e-8
X = (X - X_mean) / X_std

# Save normalization params for inference
np.savez("normalization.npz", mean=X_mean, std=X_std)
print("Saved normalization parameters")

# PyTorch Dataset
class GameDataset(Dataset):
    def __init__(self, X, y):
        self.X = torch.FloatTensor(X)
        self.y = torch.FloatTensor(y)
    
    def __len__(self):
        return len(self.X)
    
    def __getitem__(self, idx):
        return self.X[idx], self.y[idx]

# Neural Network
class GameAI(nn.Module):
    def __init__(self, input_size, output_size):
        super(GameAI, self).__init__()
        self.fc = nn.Sequential(
            nn.Linear(input_size, 128),
            nn.ReLU(),
            nn.Dropout(0.2),
            nn.Linear(128, 64),
            nn.ReLU(),
            nn.Dropout(0.2),
            nn.Linear(64, output_size),
            nn.Sigmoid()  # Multi-label output
        )
    
    def forward(self, x):
        return self.fc(x)

# Create model
input_size = X.shape[1]
output_size = y.shape[1]
model = GameAI(input_size, output_size)
print(f"Model: {input_size} -> {output_size}")

# Training
dataset = GameDataset(X, y)
loader = DataLoader(dataset, batch_size=64, shuffle=True)

criterion = nn.BCELoss()
optimizer = optim.Adam(model.parameters(), lr=0.001)

print("\nTraining...")
epochs = 50
for epoch in range(epochs):
    total_loss = 0
    for batch_X, batch_y in loader:
        optimizer.zero_grad()
        outputs = model(batch_X)
        loss = criterion(outputs, batch_y)
        loss.backward()
        optimizer.step()
        total_loss += loss.item()
    
    if (epoch + 1) % 10 == 0:
        print(f"Epoch {epoch+1}/{epochs}, Loss: {total_loss/len(loader):.4f}")

# Save model
torch.save(model.state_dict(), "ghost_hunter_ai.pth")
print("\nModel saved to ghost_hunter_ai.pth")

# Export as ONNX for Go interoperability
dummy_input = torch.zeros(1, input_size)
torch.onnx.export(model, dummy_input, "ghost_hunter_ai.onnx", 
                  input_names=['input'], 
                  output_names=['output'],
                  dynamic_axes={'input': {0: 'batch'}, 'output': {0: 'batch'}})

print("Model exported as ONNX")

# Also export as simple JSON for direct use
print("\nExporting model weights as JSON...")
state_dict = model.state_dict()
model_data = {
    "input_size": input_size,
    "output_size": output_size,
    "weights": {k: v.tolist() for k, v in state_dict.items()},
    "mean": X_mean.tolist(),
    "std": X_std.tolist()
}
with open("model_weights.json", "w") as f:
    json.dump(model_data, f)

print("Model weights exported to model_weights.json")
print("\nDone! Copy model_weights.json to the Go project to use the AI.")