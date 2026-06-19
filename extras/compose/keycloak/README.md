# Keycloak - Einrichtung für Kiosk-Go

Keycloak läuft im Dev-Modus mit eingebetteter H2-Datenbank.  
**Wichtig:** Alle Konfigurationen gehen verloren, sobald der Container mit `docker compose down` gestoppt wird. Verwende `docker compose stop` um die Daten zu erhalten.

---

## Container starten

```bash
# Aus dem Projektstamm:
docker compose -f extras/compose/keycloak/compose.yml up -d

# Stoppen (Daten bleiben erhalten):
docker compose -f extras/compose/keycloak/compose.yml stop

# Stoppen + Container löschen (Daten weg!):
docker compose -f extras/compose/keycloak/compose.yml down
```

Admin-UI: **http://localhost:8880**  
Login: `admin` / `admin`

---

## Erstkonfiguration (nach jedem `down`)

### 1. Realm anlegen

1. Links oben auf **Keycloak** (Dropdown) klicken
2. **Create realm**
3. Realm name: `kiosk`
4. **Create**

### 2. Client anlegen

1. Links → **Clients** → **Create client**
2. Client ID: `kiosk-client`
3. **Next**
4. **Direct access grants**: ON
5. **Save**

### 3. User anlegen

1. Links → **Users** → **Create new user**
2. Username: `user1`
3. Email: `user1@kiosk.de`
4. Email verified: **ON**
5. **Create**
6. Tab **Credentials** → **Set password**
7. Passwort: `p`, Temporary: **OFF** → **Save**

### 4. Verify Profile deaktivieren

1. Links → **Authentication** → **Required actions**
2. Bei **Verify Profile** → Default action: **OFF**

---

## Token holen

### Bruno

Kollektion `Kiosk-Go` → Request **POST Token holen** → Send

### curl

```bash
curl -s -X POST http://localhost:8880/realms/kiosk/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password&client_id=kiosk-client&username=user1&password=p"
```

Antwort enthält `access_token` (Bearer JWT, gültig 300 Sekunden).

---

## JWKS-Endpunkt

```
http://localhost:8880/realms/kiosk/protocol/openid-connect/certs
```

Der Endpunkt ist dokumentiert, falls später wieder eine JWT-Prüfung eingebunden
werden soll. In der aktuellen Abgabe sind die REST-Routen ungeschützt; eine
JWT-Middleware ist nicht aktiv und `internal/middleware/auth.go` ist nicht im
Repository enthalten.
