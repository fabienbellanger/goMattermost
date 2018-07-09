# Envoi de message à Mattermost

## Installation
- Installation de Go : [Golang website](https://golang.org/doc/install#install)
- Lancer
    ```
    make deps
    ```
- Dupliquer et renommer le fichier `settings.json.dist` en `settings.json` et renseigner les bonnes informations
- Exécuter le fichier `.sql` dans `database` pour initialiser la base de données (facultatif)

## Utilisation
- Envoi de commande à Mattermost avec une base de données
    ```
    make mattermost
    ```
- Envoi de commande à Mattermost sans base de données
    ```
    make mattermostNoDB
    ```
- Exécuter les tests
    ```
    make test
    ```
- Compiler et générer un binaire
    ```
    make build
    ```

## TODO
Dans toolbox :
- conversion strint to int
- conversion RawByte (date SQL) en string
