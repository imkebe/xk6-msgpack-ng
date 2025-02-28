import msgpack from 'k6/x/msgpack';

export default function() {
    try {
        // First, verify the module is loaded
        console.log("Module loaded:", msgpack);
        console.log("Functions available:", Object.keys(msgpack));
        
        // Test with a complex object containing binary data
        const testData = {
            binaryData: new Uint8Array([104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100]),
            nestedObject: {
                name: "Test Object",
                value: 123,
                timestamp: new Date().toISOString()
            },
            numbers: [1, 2, 3, 4, 5],
            isActive: true
        };
        console.log("\nTesting with complex object:");
        console.log(JSON.stringify(testData, null, 2));
        
        // Try to serialize
        console.log("Attempting to serialize...");
        const encoded = msgpack.serialize(testData);
        console.log("Serialized result type:", typeof encoded);
        //console.log("Serialized result:", encoded);
        
        // Try to deserialize
        console.log("\nAttempting to deserialize...");
        const decoded = msgpack.deserialize(encoded);
        console.log("Deserialized result type:", typeof decoded);
        console.log("Deserialized result:", decoded);
        
        // Check if round-trip worked

        
        if (testData !== decoded) {
            console.log("Values don't match!");
            console.log("Original:", testData);
            console.log("Decoded:", decoded);
        }
    } catch (error) {
        console.error("\nError occurred!");
        console.error("Error message:", error.message);
        console.error("Error type:", error.constructor.name);
        if (error.stack) {
            console.error("Stack trace:", error.stack);
        }
        throw error;
    }
}