(function() {

  function AuthDialog(titleText, actionText) {
    this._shielding = document.createElement('div');
    this._shielding.className = 'auth-dialog-shielding';
    this._dialog = document.createElement('div');
    this._dialog.className = 'auth-dialog';

    var title = document.createElement('h1');
    title.className = 'auth-dialog-title';
    title.innerText = titleText;
    this._dialog.appendChild(title);

    var lock = document.createElement('div');
    lock.className = 'auth-dialog-lock';
    this._dialog.appendChild(lock);

    this._input = document.createElement('input');
    this._input.className = 'auth-dialog-input';
    this._input.type = 'password';
    this._dialog.appendChild(this._input);

    var buttons = document.createElement('div');
    buttons.className = 'auth-dialog-buttons';
    var cancelButton = document.createElement('div');
    cancelButton.innerText = 'CANCEL';
    cancelButton.className = 'auth-dialog-button auth-dialog-button-cancel';
    buttons.appendChild(cancelButton);
    var okButton = document.createElement('div');
    okButton.innerText = actionText;
    okButton.className = 'auth-dialog-button';
    buttons.appendChild(okButton);
    this._dialog.appendChild(buttons);
  }

  AuthDialog.prototype.show = function() {
    document.body.appendChild(this._shielding);
    document.body.appendChild(this._dialog);
  };

  window.AuthDialog = AuthDialog;

})();
