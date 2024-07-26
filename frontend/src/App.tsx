import React, { useState } from "react";

interface Challenge {
  challenger: string;
  encrypted_path: string;
  encryption_method: string;
  expires_in: string;
  hint: string;
  instructions: string;
  level: number;
}

const decodeASCIIValues = (asciiStr: string): string => {
  asciiStr = asciiStr.trim().slice(1, -1); // Remove the square brackets
  const asciiValues = asciiStr
    .split(",")
    .map((val) => parseInt(val.trim(), 10));
  return String.fromCharCode(...asciiValues);
};

const App: React.FC = () => {
  const [challenge, setChallenge] = useState<Challenge | null>(null);
  const [response, setResponse] = useState<string>("");
  const [decodedPath, setDecodedPath] = useState<string>("");

  const getChallenge = async () => {
    try {
      const res = await fetch("/get-challenge");
      const data: Challenge = await res.json();
      setChallenge(data);
    } catch (error) {
      console.error("Error fetching challenge:", error);
    }
  };

  const followChallenge = async () => {
    if (challenge) {
      try {
        const res = await fetch(
          `/follow-challenge/${challenge.encrypted_path}`
        );
        const data = await res.json();
        setResponse(data);

        if (
          challenge.encryption_method ===
          "converted to a JSON array of ASCII values"
        ) {
          const decoded = decodeASCIIValues(challenge.encrypted_path);
          setDecodedPath(decoded);
        }
      } catch (error) {
        console.error("Error following challenge:", error);
      }
    }
  };

  return (
    <div>
      <h1>Challenge App</h1>
      <button onClick={getChallenge}>Get Challenge</button>
      {challenge && (
        <div>
          <h2>Challenge Details:</h2>
          <pre>{JSON.stringify(challenge, null, 2)}</pre>
          <button onClick={followChallenge}>Follow Challenge</button>
        </div>
      )}
      {response && (
        <div>
          <h2>Response:</h2>
          <pre>{response}</pre>
        </div>
      )}
      {decodedPath && (
        <div>
          <h2>Decoded Path:</h2>
          <pre>{decodedPath}</pre>
        </div>
      )}
    </div>
  );
};

export default App;
