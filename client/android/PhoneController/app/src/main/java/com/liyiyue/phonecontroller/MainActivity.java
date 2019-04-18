package com.liyiyue.phonecontroller;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.MotionEvent;
import android.widget.TextView;

import java.io.OutputStream;
import java.net.Socket;
import java.util.concurrent.ArrayBlockingQueue;

public class MainActivity extends AppCompatActivity {

	public static int FrameMilli = 1000 / 20; // 1秒20次采样
	public static int SquareLen = 10 * 10; // 阈值10像素
	public static int Scale = 16; // 像素与向量的缩放比例
	public static double N = 1.8; // 高次曲线n

	public Socket socket;
	public OutputStream os;
	public TextView text;

	public ArrayBlockingQueue<double[]> queue = new ArrayBlockingQueue<>(128);

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		text = (TextView) findViewById(R.id.text);
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_main);

		//开网络线程
		new Thread(new Runnable() {
			@Override
			public void run() {
				//connect server
				try {
					socket = new Socket("192.168.199.122", 1701);
					os = socket.getOutputStream();
				} catch (Exception e) {
					e.printStackTrace();
				}
				//start loop
				double[] lastP = new double[]{0, 0};
				long lastT = 0;
				while (true) {
					try {
						double[] vector = queue.take();
						long curT = System.currentTimeMillis();
						if ((vector[0] == 0 && vector[1] == 0) ||
								(checkPoint(lastP, vector) && checkTime(lastT, curT))) {
							//scale 1
							vector[0] *= Scale;
							vector[1] *= Scale;
							//trim
							if (vector[0] > 32700) {
								vector[0] = 32700;
							}
							if (vector[0] < -32700) {
								vector[0] = -32700;
							}
							if (vector[1] > 32700) {
								vector[1] = 32700;
							}
							if (vector[1] < -32700) {
								vector[1] = -32700;
							}
							//scale 2
							double[] sign = new double[2];
							sign[0] = vector[0] < 0 ? -1 : 1;
							sign[1] = vector[1] < 0 ? -1 : 1;
							vector[0] = Math.pow(Math.abs(vector[0]) / 32700, N) * 32700 * sign[0];
							vector[1] = Math.pow(Math.abs(vector[1]) / 32700, N) * 32700 * sign[1];
							//encode
							byte[] b = new byte[5];
							b[0] = 1; //id
							short x = (short) vector[0];
							short y = (short) vector[1];
							b[1] = (byte) (x >> 8 & 0xff);
							b[2] = (byte) (x & 0xff);
							b[3] = (byte) (y >> 8 & 0xff);
							b[4] = (byte) (y & 0xff);
							//send
							os.write(b);
							os.flush();
							//update
							lastP = vector;
							lastT = curT;
						}
					} catch (Exception e) {
						e.printStackTrace();
					}
				}
			}

			//两点距离大于阈值
			private boolean checkPoint(double[] last, double[] cur) {
				return (Math.pow(last[0] - cur[0], 2) + Math.pow(last[1] - cur[1], 2)) > SquareLen;
			}

			//两点距离大于阈值
			private boolean checkTime(long last, long cur) {
				return cur - last >= FrameMilli;
			}
		}).start();
	}

	public double startX = 0;
	public double startY = 0;

	@Override
	public boolean onTouchEvent(MotionEvent event) {
		//屏幕左上角是 {0, 0}
		try {
			switch (event.getAction()) {
				case MotionEvent.ACTION_DOWN:
					startX = event.getX();
					startY = event.getY();
					break;
				case MotionEvent.ACTION_MOVE:
					double curX = event.getX();
					double curY = event.getY();
					queue.put(new double[]{curY - startY, curX - startX});
					break;
				case MotionEvent.ACTION_UP:
					queue.put(new double[]{0, 0});
					break;
			}
		} catch (Exception e) {
			e.printStackTrace();
		}
		return super.onTouchEvent(event);
	}
}
