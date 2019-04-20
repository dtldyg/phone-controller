package com.liyiyue.phonecontroller;

import android.content.Context;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.DisplayMetrics;
import android.view.MotionEvent;
import android.view.WindowManager;
import android.widget.TextView;

import java.io.IOException;
import java.io.OutputStream;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;
import java.util.concurrent.ArrayBlockingQueue;

public class MainActivity extends AppCompatActivity {

	public static int MilliPerFrame = 1000 / 30; // 1秒n次采样
	public static int SquareLen = 1; // 阈值n像素以内认为是一个点
	public static double Scale = 1; // 像素与向量的缩放比例

	public Socket socketTcp;
	public OutputStream osTcp;
	public DatagramSocket socketUdp;
	public SocketAddress addressUdp;
	public TextView text;

	public int width;
	public int heigh;
	public int scrollBarWidth;

	public ArrayBlockingQueue<Byte> queueAction = new ArrayBlockingQueue<>(128);
	public ArrayBlockingQueue<double[]> queueStatus = new ArrayBlockingQueue<>(128);

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_main);
		text = (TextView) findViewById(R.id.text);

		width = getScreenWidth(getApplicationContext());
		heigh = getScreenHeight(getApplicationContext());
		scrollBarWidth = getResources().getDimensionPixelSize(R.dimen.scroll_bar);

		//开tcp线程
		new Thread(new Runnable() {
			@Override
			public void run() {
				//connect server
				try {
					socketTcp = new Socket("192.168.199.122", 1701);
					osTcp = socketTcp.getOutputStream();
				} catch (Exception e) {
					e.printStackTrace();
				}
				//start loop
				while (true) {
					try {
						byte action = queueAction.take();
						//build
						byte[] b = new byte[2];
						b[0] = 2; //id
						b[1] = action;
						osTcp.write(b);
					} catch (InterruptedException e) {
						e.printStackTrace();
					} catch (IOException e) {
						e.printStackTrace();
					}
				}
			}
		}).start();

		//开udp线程
		new Thread(new Runnable() {
			@Override
			public void run() {
				//connect server
				try {
					socketUdp = new DatagramSocket();
					addressUdp = new InetSocketAddress("192.168.199.122", 1702);
				} catch (Exception e) {
					e.printStackTrace();
				}
				//start loop
				double[] lastP = new double[]{0, 0};
				long lastT = 0;
				while (true) {
					try {
						double[] curP = queueStatus.take();
						long curT = System.currentTimeMillis();
						if ((lastP[0] == 0 && lastP[1] == 0) || (curP[0] == 0 && curP[1] == 0)) {
							lastP = curP;
							continue;
						}
						short dx = (short) ((curP[0] - lastP[0]) * Scale);
						short dy = (short) ((curP[1] - lastP[1]) * Scale);
						if (checkPoint(dx, dy) && checkTime(lastT, curT)) {
							//build
							byte[] b = new byte[5];
							b[0] = 1; //id
							short x = dx;
							short y = dy;
							b[1] = (byte) (x >> 8 & 0xff);
							b[2] = (byte) (x & 0xff);
							b[3] = (byte) (y >> 8 & 0xff);
							b[4] = (byte) (y & 0xff);
							//send
							DatagramPacket pak = new DatagramPacket(b, 5, addressUdp);
							socketUdp.send(pak);
							//reset
							lastP[0] += dx;
							lastP[1] += dy;
						}
					} catch (InterruptedException e) {
						e.printStackTrace();
					} catch (IOException e) {
						e.printStackTrace();
					}
				}
			}

			//两点距离大于阈值
			private boolean checkPoint(short dx, short dy) {
				return (Math.pow(dx, 2) + Math.pow(dy, 2)) >= SquareLen;
			}

			//采样间隔大于阈值
			private boolean checkTime(long last, long cur) {
				return cur - last >= MilliPerFrame;
			}
		}).start();
	}

	int leftKey = -1;
	int rightKey = -1;
	int moveKey = -1;

	@Override
	public boolean onTouchEvent(MotionEvent event) {
		try {
			int id = event.getPointerId(event.getActionIndex());
			switch (event.getActionMasked()) {
				case MotionEvent.ACTION_DOWN:
					down(event, id);
					break;
				case MotionEvent.ACTION_POINTER_DOWN:
					down(event, id);
					break;
				case MotionEvent.ACTION_MOVE:
					move(event, id);
					break;
				case MotionEvent.ACTION_UP:
					up(event, id);
					break;
				case MotionEvent.ACTION_POINTER_UP:
					up(event, id);
					break;
			}
		} catch (InterruptedException e) {
			e.printStackTrace();
		}
		return true;
	}

	public void down(MotionEvent event, int id) throws InterruptedException {
		double x = event.getX(id);
		double y = event.getY(id);
		System.out.println(x + "," + y);
		System.out.println(width + "," + heigh + "," + scrollBarWidth);
		if (x > width / 2 && x < width - scrollBarWidth && y < heigh / 2 && leftKey == -1) {
			//左键按下
			leftKey = id;
			queueAction.put((byte) 1);
		} else if (x > width / 2 && x < width - scrollBarWidth && y > heigh / 2 && rightKey == -1) {
			//右键按下
			rightKey = id;
			queueAction.put((byte) 3);
		} else if (x < width / 2 && moveKey == -1) {
			moveKey = id;
		}
	}

	public void move(MotionEvent event, int id) throws InterruptedException {
		if (moveKey > -1) {
			queueStatus.put(new double[]{event.getX(moveKey), event.getY(moveKey)});
		}
	}

	public void up(MotionEvent event, int id) throws InterruptedException {
		if (moveKey > -1 && id == moveKey) {
			moveKey = -1;
			queueStatus.put(new double[]{0, 0});
		} else if (leftKey > -1 && id == leftKey) {
			//左键释放
			leftKey = -1;
			queueAction.put((byte) 2);
		} else if (rightKey > -1 && id == rightKey) {
			//右键释放
			rightKey = -1;
			queueAction.put((byte) 4);
		}
	}

	/**
	 * 获取屏幕的宽
	 *
	 * @param context
	 * @return
	 */
	public static int getScreenWidth(Context context) {
		WindowManager wm = (WindowManager) context.getSystemService(Context.WINDOW_SERVICE);
		DisplayMetrics dm = new DisplayMetrics();
		wm.getDefaultDisplay().getMetrics(dm);
		return dm.widthPixels;
	}

	/**
	 * 获取屏幕的高度
	 *
	 * @param context
	 * @return
	 */
	public static int getScreenHeight(Context context) {
		WindowManager wm = (WindowManager) context.getSystemService(Context.WINDOW_SERVICE);
		DisplayMetrics dm = new DisplayMetrics();
		wm.getDefaultDisplay().getMetrics(dm);
		return dm.heightPixels;
	}
}
