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
- Envoi de commande à Mattermost et/ou Slack avec une base de données
    ```
    make notification
    ```
- Envoi de commande à Mattermost et/ou Slack sans base de données
    ```
    make notificationNoDB
    ```
- Initialiser la base de données (dry-run)
    ```
    make db
    ```
- Initialiser la base de données (force)
    ```
    make dbForce
    ```
- Exécuter les tests
    ```
    make test
    ```
- Compiler et générer un binaire
    ```
    make build
    ```
- Lancement du serveur Web
    ```
    make serve
    ```
- Envoi du mail quotidien
    ```
    make mail
    ```

## TODO
Dans toolbox :
- conversion string <=> int
- conversion RawByte (date SQL) en string

## Envoi de mails avec Maildev

Lancer : `maildev --incoming-user 'root' --incoming-pass 'root'`
