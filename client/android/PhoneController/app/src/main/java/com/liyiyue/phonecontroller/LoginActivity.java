package com.liyiyue.phonecontroller;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

public class LoginActivity extends AppCompatActivity {

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_login);

		Button btn = (Button) findViewById(R.id.btn_login);
		btn.setOnClickListener(new View.OnClickListener() {
			@Override
			public void onClick(View v) {
				EditText etIP1 = (EditText) findViewById(R.id.ip_1);
				EditText etIP2 = (EditText) findViewById(R.id.ip_2);
				EditText etIP3 = (EditText) findViewById(R.id.ip_3);
				EditText etIP4 = (EditText) findViewById(R.id.ip_4);
				String ip = etIP1.getText() + "." + etIP2.getText() + "." + etIP3.getText() + "." + etIP4.getText();

				Intent intent = new Intent(LoginActivity.this, MainActivity.class);
				intent.putExtra("ip", ip);
				startActivity(intent);
			}
		});
	}
}
