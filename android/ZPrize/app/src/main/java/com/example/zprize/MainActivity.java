package com.example.zprize;

import android.os.Build;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import android.text.InputType;
import android.text.TextUtils;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.TextView;

import androidx.annotation.RequiresApi;
import androidx.appcompat.app.AppCompatActivity;

import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;

@RequiresApi(api = Build.VERSION_CODES.P)
public class MainActivity extends AppCompatActivity {
    private final Handler mainThreadHandler = Handler.createAsync(Looper.getMainLooper());
    private static final String TAG = "ZPrizeTestHarness";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        LinearLayout linearLayout = findViewById(R.id.rootContainer);

        Button btnShow = new Button(this);
        btnShow.setText("Press to run using random elements");
        btnShow.setLayoutParams(new LinearLayout.LayoutParams(ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT));

        Button btnShowFile = new Button(this);
        btnShowFile.setText("Press to run from test vector file");
        btnShowFile.setLayoutParams(new LinearLayout.LayoutParams(ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT));

        EditText iters = new EditText(this);
        iters.setHint("#iterations per vector");
        iters.setInputType(InputType.TYPE_CLASS_NUMBER);

        EditText numElems = new EditText(this);
        numElems.setHint("#elems as power of 2");
        numElems.setInputType(InputType.TYPE_CLASS_NUMBER);

        EditText numVecs = new EditText(this);
        numVecs.setHint("#vectors to generate randomly");
        numVecs.setInputType(InputType.TYPE_CLASS_NUMBER);

        TextView resultView = new TextView(this);
        TextView resultView2 = new TextView(this);

        File filePoints = new File(getFilesDir()+"/points");
        try {
            InputStream is = getAssets().open("points");
            FileOutputStream fos = new FileOutputStream(filePoints);
            int size = Math.min(is.available(), 1 << 14);
            byte[] buffer = new byte[size];
            while (true) {
                int len = is.read(buffer);
                if (len == -1) { break; }
                fos.write(buffer, 0, len);
            }
            is.close();
            fos.close();
        } catch (Exception e) { throw new RuntimeException(e); }

        File fileScalars = new File(getFilesDir()+"/scalars");
        try {
            InputStream is = getAssets().open("scalars");
            FileOutputStream fos = new FileOutputStream(fileScalars);
            int size = Math.min(is.available(), 1 << 14);
            byte[] buffer = new byte[size];
            while (true) {
                int len = is.read(buffer);
                if (len == -1) { break; }
                fos.write(buffer, 0, len);
            }
            is.close();
            fos.close();
        } catch (Exception e) { throw new RuntimeException(e); }

        btnShow.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                resultView.setText("Running on random vectors");
                resultView2.setText("");
                File dir = getFilesDir();
                String dir_path = dir.getAbsolutePath();
                String iters_val = iters.getText().toString();
                String numElemsVal = numElems.getText().toString();
                String numVecsVal = numVecs.getText().toString();
                if (TextUtils.isDigitsOnly(iters_val) && !TextUtils.isEmpty(iters_val)
                && TextUtils.isDigitsOnly(numElemsVal) && !TextUtils.isEmpty(numElemsVal)) {
                    new Thread(new Runnable() {
                        @Override
                        public void run() {
                            RustMSM g = new RustMSM();
                            Log.i(TAG, "Starting MSM with random inputs");
                            String r = g.runMSMRandomMultipleVecs(dir_path, iters_val, numElemsVal, numVecsVal);
                            Log.i(TAG, "Finished MSM with random inputs");

                            // doing the same thing with gnark
                            GnarkMSM gnarkMSM = new GnarkMSM();
                            String rGnark = gnarkMSM.runMSMRandomMultipleVecs(dir_path, iters_val, numElemsVal, numVecsVal);


                            mainThreadHandler.post(new Runnable() {
                                @Override
                                public void run() {
                                    String result = "Mean time to run with random elements is: ";
                                    resultView.setText(result);
                                    resultView2.setText("[ ref ]: " + r + "\n[gnark]: " + rGnark);
                                }
                            });
                        }
                    }).start();
                } else {
                    resultView.setText("Valid number of iterations, vectors, and elements per vector must be provided");
                    resultView2.setText("");
                }
            }
        });

        btnShowFile.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                File dir = getFilesDir();
                String dir_path = dir.getAbsolutePath();
                String iters_val = iters.getText().toString();
                if (TextUtils.isDigitsOnly(iters_val) && !TextUtils.isEmpty(iters_val)) {
                    resultView.setText("Currently running test vectors");
                    resultView2.setText("");

                    new Thread(new Runnable() {
                        @Override
                        public void run() {
                            RustMSM g = new RustMSM();
                            Log.i(TAG, "Starting MSM with fixed inputs");
                            String r =  g.runMSMFile(dir_path, iters_val);
                            Log.i(TAG, "Finished MSM with fixed inputs");

                            GnarkMSM gnarkMSM = new GnarkMSM();
                            String rGnark = gnarkMSM.runMSMFile(dir_path, iters_val);


                            mainThreadHandler.post(new Runnable() {
                                @Override
                                public void run() {
                                    String result = "Mean time to run with test vectors is: ";
                                    resultView.setText(result);
                                    resultView2.setText("[ ref ]: " + r + "\n[gnark]: " + rGnark);
                                }
                            });
                        }
                    }).start();
                } else {
                    resultView.setText("Valid number of iterations must be provided");
                    resultView2.setText("");
                }
            }
        });

        // Add Button to LinearLayout
        if (linearLayout != null) {
            linearLayout.addView(btnShowFile);
            linearLayout.addView(btnShow);
            linearLayout.addView(iters);
            linearLayout.addView(numVecs);
            linearLayout.addView(numElems);
            linearLayout.addView(resultView);
            linearLayout.addView(resultView2);
        }
    }

    static {
        System.loadLibrary("msm");
    }
}