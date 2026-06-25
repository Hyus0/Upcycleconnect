import mysql.connector
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.preprocessing import LabelEncoder
import numpy as np

conn = mysql.connector.connect(
    host="localhost",
    user="root",
    password="",
    database="new_upcycle" 
)
cursor = conn.cursor()

cursor.execute("""
CREATE TABLE IF NOT EXISTS ML_PREDICTION (
    id_utilisateur INT PRIMARY KEY,
    prestation_recommandee VARCHAR(100),
    confiance DECIMAL(5,2),
    date_prediction DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEUR(id) ON DELETE CASCADE
)
""")
cursor.execute("TRUNCATE TABLE ML_PREDICTION") 

print("Récupération des données depuis SQL...")
query = """
    SELECT u.id, u.role, u.ville, 
           COALESCE(COUNT(a.id), 0) as nb_annonces
    FROM UTILISATEUR u
    LEFT JOIN ANNONCE a ON u.id = a.id_vendeur
    GROUP BY u.id
"""
df = pd.read_sql(query, conn)

def determine_target(row):
    if row['role'] == 'Prestataire': return 'Abonnement Pro'
    elif row['nb_annonces'] > 2: return 'Dépôt Conteneur'
    elif row['ville'] == 'Paris': return 'Formation Upcycling'
    else: return 'Achat Matériaux'

df['prestation_cible'] = df.apply(determine_target, axis=1)

print("Entraînement du modèle de classification...")
le_role = LabelEncoder()
le_ville = LabelEncoder()

X = pd.DataFrame()
X['role_encoded'] = le_role.fit_transform(df['role'])
X['ville_encoded'] = le_ville.fit_transform(df['ville'])
X['nb_annonces'] = df['nb_annonces']
y = df['prestation_cible']

clf = RandomForestClassifier(n_estimators=100, random_state=42)
clf.fit(X, y)

predictions = clf.predict(X)
probabilites = clf.predict_proba(X)
confiances = np.max(probabilites, axis=1) * 100

print("Sauvegarde des prédictions dans SQL...")
insert_query = """
    INSERT INTO ML_PREDICTION (id_utilisateur, prestation_recommandee, confiance)
    VALUES (%s, %s, %s)
"""

data_to_insert = []
for i in range(len(df)):
    data_to_insert.append((
        int(df.iloc[i]['id']),
        str(predictions[i]),
        float(confiances[i])
    ))

cursor.executemany(insert_query, data_to_insert)
conn.commit()

print(f"Terminé ! {len(data_to_insert)} prédictions insérées dans la base de données.")
cursor.close()
conn.close()