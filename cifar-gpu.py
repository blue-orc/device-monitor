import sys, getopt
import numpy as np
import pickle
import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
import torchvision
import torchvision.transforms as transforms
from multiprocessing import Process, Manager

# Reference documentation:
# https://pytorch.org/tutorials/beginner/blitz/cifar10_tutorial.html

layers = 2
depth = 300
batchSize = 2000
numWorkers = 4
epochs = 10
numberOfFiles = 10
learningRate = 0.001

def writeOutput(key, value):
    separator = ":"
    print(key + separator + str(value))
    sys.stdout.flush()

base = "/mnt/cifar-data/data_batch_"

# Dataset has 1,000 files each containing 10,000 images

# Deserializes the data
def unpickle(file):
    with open(file, 'rb') as fo:
        dict = pickle.load(fo, encoding='bytes')
    return dict

def getDataset(D, x):
    t = MyDataset(base+ str(x), transform=transform)
    trainloader = torch.utils.data.DataLoader(t, batch_size=batchSize,
                                          shuffle=True, num_workers=numWorkers)
    D.append(trainloader)

# Define the convolutional neural network (CNN)
class Net(nn.Module):
    def __init__(self):
        super(Net, self).__init__()
        self.conv1 = nn.Conv2d(3, depth, 5)
        self.pool = nn.MaxPool2d(2, 2)
        self.conv2 = nn.Conv2d(depth, 16, 5)
        self.fc1 = nn.Linear(16 * 5 * 5, 120)
        self.fc2 = nn.Linear(120, 84)
        self.fc3 = nn.Linear(84, 10)

    def forward(self, x):
        x = self.pool(F.relu(self.conv1(x)))
        x = self.pool(F.relu(self.conv2(x)))
        x = x.view(-1, 16 * 5 * 5)
        x = F.relu(self.fc1(x))
        x = F.relu(self.fc2(x))
        x = self.fc3(x)
        return x

# Define the dataset
# Controls how PyTorch dataloader will access the data
class MyDataset(torch.utils.data.Dataset):
    def __init__(self, file, transform=None):
        ds = unpickle(file)
        data = np.asarray(ds[b'data'])
        shaped = np.reshape(data, (10000,32,32,3))
        self.data = shaped
        self.target = ds[b'labels']
        self.transform = transform
        
    def __getitem__(self, index):
        x = self.data[index]
        y = self.target[index]

        if self.transform:
            x = self.transform(x)
        
        return x, y
    
    def __len__(self):
        return len(self.target)

transform = transforms.Compose(
    [transforms.ToTensor(),
     transforms.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))])


writeOutput("Status", "Running")
writeOutput("Epochs", epochs)
writeOutput("NumberOfWorkers", numWorkers)
writeOutput("NumberOfFiles", numberOfFiles)
writeOutput("BatchSize", batchSize)
writeOutput("Depth", depth)
writeOutput("Layers", layers)
writeOutput("LearningRate", learningRate)
writeOutput("TrainingScript", "cifar-gpu.py")
device = torch.device("cuda:0")
print(device)

writeOutput("Step", "Creating Neural Network")
net = Net()
net = net.to(device)
criterion = nn.CrossEntropyLoss()
optimizer = optim.SGD(net.parameters(), lr=learningRate, momentum=0.9)

manager = Manager()
D = manager.list()  # <-- can be shared between processes.
processes = []
writeOutput("Step", "Reading Image Files")
for i in range(1,numberOfFiles + 2):
    p = Process(target=getDataset, args=(D,i))  # Passing the list
    p.start()
    processes.append(p)
for p in processes:
    p.join()

writeOutput("Step", "Training")
writeOutput("NumberOfFiles", len(D))
writeOutput("ImagesPerFile", len(D[0]))

# Loop over the dataset multiple times. Change this value to run container longer.
for epoch in range(epochs):  
    running_loss = 0.0
    # Every epoch, load each file in dataset and iterate over the data
    for x in range(1,len(D)):
        
        for i, data in enumerate(D[x], 0):
            # Get the inputs; data is a list of [inputs, labels]
            inputs, labels = data[0].to(device), data[1].to(device)

            # Zero the parameter gradients
            optimizer.zero_grad()

            # forward + backward + optimize
            outputs = net(inputs)
            loss = criterion(outputs, labels)
            loss.backward()
            optimizer.step()

            # Print statistics
            running_loss += loss.item()
            if i % 1 == 0:
                writeOutput("CurrentFileIndex", x)
                writeOutput("CurrentImageIndex", (i+1)*2000)
                writeOutput("CurrentEpoch", epoch+1)
                writeOutput("Loss", running_loss / 2000)
                running_loss = 0.0

writeOutput("Step", "Training Complete")
writeOutput("Status", "Finished")