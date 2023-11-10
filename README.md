### Connects to neo4j

### Create data as

`


    CREATE

    
    (node12:NODE {id: "m12", name: "m12"}),

    
    (node13:NODE {id: "m13", name: "m13"}),

    
    (node14:NODE {id: "m14", name: "m14"}),

    
    (node15:NODE {id: "m15", name: "m15"}),

    
    (node16:NODE {id: "m16", name: "m16"}),

    
    (node17:NODE {id: "m17", name: "m17"}),

    
    (node18:NODE {id: "m18", name: "m18"}),

    
    (node19:NODE {id: "m19", name: "m19"}),

    
    (node20:NODE {id: "m20", name: "m20"}),

    
    (node21:NODE {id: "m21", name: "m21"}),

    
    (node12)-[:NEXT {w:1, optional: false}]->(node13),

    
    (node13)-[:NEXT {w:1, optional: false}]->(node14),

    
    (node13)-[:NEXT {w:2, optional: false}]->(node15),

    
    (node13)-[:NEXT {w:3, optional: false}]->(node16),

    
    (node14)-[:NEXT {w:1, optional: false}]->(node17),

    
    (node15)-[:NEXT {w:1, optional: false}]->(node17),

    
    (node16)-[:NEXT {w:1, optional: false}]->(node17),

    
    (node17)-[:NEXT {w:1, optional: false}]->(node18),

    
    (node18)-[:NEXT {w:1, optional: false}]->(node19),

    
    (node18)-[:NEXT {w:2, optional: false}]->(node20),

    
    (node18)-[:NEXT {w:3, optional: false}]->(node21);
`

